// Copyright 2022 gab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var DEFAULT_GAME_TRAFFIC_FREQ = 1
var DEFAULT_GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ = 30

type Config struct {
	Port                            string
	DbURI                           string
	GameTrafficFreq                 time.Duration
	GameTrafficLimitPlayTimePerFreq int32
}

func EnvConfig() Config {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		panic("PORT not provided")
	}

	dbURI, ok := os.LookupEnv("POSTGRESQL_URL")

	if !ok {
		panic("POSTGRESQL_URL not provided")
	}

	var gameTrafficFreq int
	gameTrafficFreqStr, ok := os.LookupEnv("GAME_TRAFFIC_FREQUENCY")

	if !ok {
		log.Print(fmt.Sprintf("GAME_TRAFFIC_FREQUENCY not provided. Default is %d min", DEFAULT_GAME_TRAFFIC_FREQ))
		gameTrafficFreq = DEFAULT_GAME_TRAFFIC_FREQ
	}

	gameTrafficFreq, err := strconv.Atoi(gameTrafficFreqStr)
	if err != nil {
		panic("GAME_TRAFFIC_FREQUENCY is not an integer")
	}

	var gameTrafficLimitPlayTimePerFreq int32
	gameTrafficLimitPlayTimePerFreqStr, ok := os.LookupEnv("GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ")

	if !ok {
		log.Print(fmt.Sprintf("GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ not provided. Default is %d", DEFAULT_GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ))
		gameTrafficLimitPlayTimePerFreq = int32(DEFAULT_GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ)
	}

	gameTrafficLimitPlayTimePerFreqInt, err := strconv.Atoi(gameTrafficLimitPlayTimePerFreqStr)
	if err != nil {
		panic("GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ is not an integer")
	}
	if gameTrafficLimitPlayTimePerFreqInt < 0 {
		panic("GAME_TRAFFIC_LIMIT_PLAY_TIME_PER_FREQ is 0, so there will be no traffic")
	}
	gameTrafficLimitPlayTimePerFreq = int32(gameTrafficLimitPlayTimePerFreqInt)

	return Config{port, dbURI, time.Duration(int32(gameTrafficFreq)) * time.Second, gameTrafficLimitPlayTimePerFreq}
}
