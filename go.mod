module timeMate_bot

go 1.19

require github.com/mattn/go-sqlite3 v1.14.22

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-00010101000000-000000000000

//require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => ./telegram-bot-api/

//replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => github.com/OvyFlash/telegram-bot-api/v5 v5.0.0-20240108230938-63e5c59035bf
