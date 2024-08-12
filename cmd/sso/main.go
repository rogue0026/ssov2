package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/rogue0026/ssov2/internal/service"
	"github.com/rogue0026/ssov2/internal/ssoconfig"
	"github.com/rogue0026/ssov2/internal/storage/postgres"
)

const (
	ENV_DEV  = "dev"
	ENV_PROD = "prod"
)

func main() {
	// loading configuration
	ssoConfigPath := flag.String("c", "", "path to sso config file")
	flag.Parse()
	appConfig := ssoconfig.MustLoad(*ssoConfigPath)

	// setup logger
	appLogger := setupLogger(appConfig.RunningEnv)

	// initializing connection to database
	s, err := postgres.New(context.Background(), appConfig.DSN)
	if err != nil {
		appLogger.Error(err.Error())
	}

	// initializing service layer
	svc := service.New(s, s)
	userId, err := svc.RegisterNewUser(context.Background(), "test_login", "test_password", "test@example.com")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("user successfully created. user id is %d\n", userId)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case ENV_DEV:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	case ENV_PROD:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))
	}

	return logger
}

/*
План работы:
0. Создать логгер (done)
1. Сделать загрузку конфигурации для работы приложения: (done)
	Среда выполнения приложения
	Адрес сервера
	Строка подключения к базе данных
	Время жизни токена
2. Написать слой для взаимодействия с базой данных PostgreSQL: (done)
	Создание (регистрация) нового пользователя в системе
	Получение данных о пользователе по Login && Email
	Удаление пользователя из системы
3. Написать сервисный слой, который будет взаимодействовать со слоем хранения данных пользователей
	Создание пользователя в системе: (done)
		Проверка входных данных для регистрации пользователя
		Создание нового пользователя или информирование клиента об ошибке создания пользователя
	Генерация нового токена для конкретного пользователя
		Проверка полученных данных об учетной записи пользователя (todo)
		Генерация нового токена в случае успеха или возврат ошибки (done)
4. Написать слой для обмена данными через gRPC
	Создать gRPC-обработчики
5. Написать gRCP-interceptors
	Логирование запросов и ответов
6. Написать слой для работы с Kafka. Через нее будет осуществляться централизованный сбор логов и метрик микросервисов для визуализации с помощью ELK и Grafana
7. Написать Dockerfile для сборки образа приложения
8. Написать docker-compose для развертывания приложения и базы данных в облаке
*/
