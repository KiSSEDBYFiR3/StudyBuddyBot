package telegram

const msgHelp = `**Доступные команды:**

/dick - узнай всё про свой 🍌
/duel _@username_ - вызвать на (или принять) бой ⚔️
/top\_dick - статистика всех 🍆

/gay - узнать у кого сегодня удачный день 🤡 (_работает, только на админов чата_)
/top\_gay - статистика по бедолагам в чате 🔞

/schedule - получить расписание из Google Calendar (_рабоает только если привязан calendar-id группы_)
/add\_calendar *[calendar-id]* - привязать расписание из Google Calendar (_возможно только для админов группы_)
(*ВАЖНО*  - _не забудьте в настройках календаря открыть доступ пользователю: calendar-manager@flash-spark-404006.iam.gserviceaccount.com_)

/flip - подбросить монетку 🪙
/xkcd - случайный xkcd комикс 😂
`

const msgHello = "Hi there\n\n" + msgHelp

const (
	msgCreateUser     = "@%s, только что обнаружил(а) свой пенис 🤣\n"
	msgAlreadyPlays   = "@%s, ты уже играл сегодня 🚨\n"
	msgDickSize       = "Теперь размер его пениса: %d см 🍌"
	msgDickIncrease   = "@%s, твой пенис увеличился на %d см 😍\n"
	msgDickDecrease   = "@%s, твой пенис уменьшился на %d см 😭\n"
	msgChangeDickSize = "@%s, %d ➜ %d см 🍌\n"

	msgTargetNotFound = "@%s этот пользователь не имеет писюна 🍆"
	msgVictoryInDuel  = "@%s победил в дуели @%s\n"
	msgUserHasBanned  = "@%s получает бан на %d секунд 🚫\n"

	msgChanceDuel = "@%s имеет пенис %d см и шансы на победу %.2f%%\n@%s имеет пенис %d см и шансы на победу %.2f%%\n"

	msgNewGayOfDay     = "Новый пидор дня - @%s"
	msgCurrentGayOfDay = "Текущий пидор дня - @%s"

	msgDuelWithYourself = "@%s засунул пенис себе в рот 🍆"

	msgChallengeToDuel = "@%s вызывает на дуель @%s"
	msgAcceptDuel      = "@%s принимает вызов на дуель (размер пениса %d см, шансы на победу %.2f%%) от @%s (размер пениса %d см, шансы на победу %.2f%%)\n"
	msgUser1Wins       = "@%s побеждает в дуели, теперь размер его пениса %d см\nРазмер пениса @%s уменьшается до %d см"
	msgUser1Lost       = "@%s проигрывает в дуели, теперь размер его пениса %d см\nРазмер пениса @%s увеличивается до %d см"

	msgForbiddenCalendarUpdate = "Обновить id календаря может только администратор группы"
	msgErrorSendMessage        = "Не удалось получить расписание из календаря с id: %s"
	msgErrorUpdateCalendarID   = "Не удалось добавить данное id календаря: \"%s\""
	msgSuccessUpdateCalendarID = "Теперь к вашему чату привязано новое расписание\nНе забудь открыть доступ к своему календарю пользователю:\ncalendar-manager@flash-spark-404006.iam.gserviceaccount.com\n"
)
