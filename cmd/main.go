package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	todo "github.com/kirD2287/REST-GO"
	handler "github.com/kirD2287/REST-GO/pkg/handler"
	"github.com/kirD2287/REST-GO/pkg/repository"
	"github.com/kirD2287/REST-GO/pkg/service"
	"github.com/spf13/viper"
)



func main() {
	if err := initConfig(); err!= nil {
        logrus.Fatalf("error reading config file: %s", err.Error())
    }

		if err := godotenv.Load() ; err!= nil {
				log.Fatalf("error loading env_variables: %s", err.Error())
		}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
        User: viper.GetString("db.user"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})
	if err!= nil {
        logrus.Fatalf("error connecting to database: %s", err.Error())
    }
	
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	
	srv := new(todo.Server)
	go func () {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

		logrus.Print("TodoApp Started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		logrus.Print("TodoApp Finished")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Errorf("error occured on server shutdown: %s", err.Error())
		}

		if err := db.Close(); err != nil {
			logrus.Errorf("error occured on db connection close: %s", err.Error())
		}
	}
	



func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}