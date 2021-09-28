package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/terdia/todo/proto/todo"
)

func main() {

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stdout, "missing subcommand: list or add")
		os.Exit(1)
	}

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to :8888: %v", err)
	}

	client := todo.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

}

func list(ctx context.Context, client todo.TasksClient) error {

	l, err := client.List(ctx, new(emptypb.Empty))
	if err != nil {
		return fmt.Errorf("could not fetch task: %v", err)
	}

	for _, task := range l.Tasks {
		if task.Done {
			fmt.Printf("✅")
		} else {
			fmt.Printf("❌")
		}

		fmt.Printf(" %s\n", task.Text)
	}
	return nil
}

func add(ctx context.Context, client todo.TasksClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{Text: text})
	if err != nil {
		return fmt.Errorf("could not add task in the backend: %v", err)
	}

	fmt.Println("task addded successfuly")

	return nil
}
