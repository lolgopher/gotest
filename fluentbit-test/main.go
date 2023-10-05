package main

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	hostname    = os.Getenv("HOSTNAME")
	programVer  = "dev"
	programName = "logurs-test"

	log zerolog.Logger
)

func main() {
	if len(hostname) == 0 {
		hostname = "local"
	}

	initLogger()

	defer func() {
		recover()
		log.Fatal().Msg("this log is fatal test log")
	}()

	log.Trace().Msg("this log is trace test log")
	log.Debug().Msg("this log is debug test log")
	log.Info().Msg("this log is info test log")
	log.Warn().Msg("this log is warn test log")
	log.Error().Msg("this log is error test log")
	log.Panic().Msg("this log is panic test log")

	// [2023-09-25T14:06:57+09:00] DEBUG this log is debug test log hostname=local program=logurs-test ver=dev
	// [2023-09-25T14:06:57+09:00]  INFO this log is info test log hostname=local program=logurs-test ver=dev
	// [2023-09-25T14:06:57+09:00]  WARN this log is warn test log hostname=local program=logurs-test ver=dev
	// [2023-09-25T14:06:57+09:00] ERROR this log is error test log hostname=local program=logurs-test ver=dev
	// [2023-09-25T14:06:57+09:00] PANIC this log is panic test log hostname=local program=logurs-test ver=dev
	// [2023-09-25T14:06:57+09:00] FATAL this log is fatal test log hostname=local program=logurs-test ver=dev

	// app.neo: [1695625348.545773323, {
	// 	"cid"=>"NX_000A0CD9651130630284_9d8a5998486f0cb13c41ba5702bfe3e2df9a5130847039e6c0-0017-8484_AI",
	//	"svc"=>"kep_ccaas-ccaas", "program"=>"neo", "ver"=>"v1.12.20-88f7751(2023-09-07_01:59:23_UTC)",
	//	"level"=>"info", "message"=>"call terminated", "sessid"=>"e0212aef-83c8-4f7d-9fee-f609ae988e44",
	//	"hostname"=>"nebuchadnezzar-ccass-callgw-1"
	// }]
}

func initLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log = zerolog.New(os.Stdout).With().Timestamp().
		Str("program", programName).
		Str("ver", programVer).
		Str("hostname", hostname).Logger()

}
