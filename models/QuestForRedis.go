package models

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)
type QuestRedis struct {
	Category	string	`redis:"category"`
	Title		string	`redis:"title"`
	Likes		int		`redis:"likes"`
	Tasks		string	`redis:"tasks"`
}
// Declare a pool variable to hold the pool of Redis connections.
var Pool *redis.Pool

var ErrNoQuest = errors.New("quest not found")
var ErrNoCategory = errors.New("category not found")

func FindQuest(id string) (*QuestRedis, error) {
	conn := Pool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", "QuestRedis:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoQuest
	}

	var quest QuestRedis
	err = redis.ScanStruct(values, &quest)
	if err != nil {
		return nil, err
	}

	return &quest, nil
}

func IncrementLikes(id string) error {
	conn := Pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", "QuestRedis:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoQuest
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("HINCRBY", "QuestRedis:"+id, "likes", 1)
	if err != nil {
		return err
	}
	// And we do the same with the increment on our sorted set.
	err = conn.Send("ZINCRBY", "likes", 1, id)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func AddRecord(questRedis QuestRedis) error {
	conn := Pool.Get()
	defer conn.Close()

	lastFreeId, err := redis.String(conn.Do("GET", "lastFreeId"))
	if err != nil {
		return err
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("INCR","lastFreeId")
	if err != nil {
		return err
	}
	err = conn.Send("HMSET",
		"QuestRedis:"+lastFreeId,
		"title",questRedis.Title,
		"likes",0,
		"tasks","do anything")
	if err != nil {
		return err
	}
	// And we do the same with the increment on our sorted set.
	err = conn.Send("ZINCRBY", "likes", 0, lastFreeId)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}


func FindTop(numbersOfTop int) ([]*QuestRedis, error) {
	conn := Pool.Get()
	defer conn.Close()
	//initial number of

	for {
		// Instruct Redis to watch the likes sorted set for any changes.
		_, err := conn.Do("WATCH", "likes")
		if err != nil {
			return nil, err
		}

		ids, err := redis.Strings(conn.Do("ZREVRANGE", "likes", 0,numbersOfTop-1))
		if err != nil {
			return nil, err
		}

		err = conn.Send("MULTI")
		if err != nil {
			return nil, err
		}


		for _, id := range ids {
			err := conn.Send("HGETALL", "QuestRedis:"+id)
			if err != nil {
				return nil, err
			}
		}

		replies, err := redis.Values(conn.Do("EXEC"))
		if err == redis.ErrNil {
			log.Println("trying again")
			continue
		} else if err != nil {
			return nil, err
		}

		quests := make([]*QuestRedis, numbersOfTop)

		for i, reply := range replies {
			var quest QuestRedis
			err = redis.ScanStruct(reply.([]interface{}), &quest)
			if err != nil {
				return nil, err
			}

			quests[i] = &quest
		}

		return quests, nil
	}
}
