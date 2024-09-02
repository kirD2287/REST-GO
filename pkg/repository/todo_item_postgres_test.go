package repository

import (
	"errors"
	"log"
	"testing"

	todo "github.com/kirD2287/REST-GO"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)




func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err!= nil {
        log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item todo.TodoItem
	}

	type mockBehavior func(args args, id int) 

	testTable := []struct {
		name string
        args args
        mockBehavior mockBehavior
		id int
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
                item: todo.TodoItem{
					Title: "test title", 
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
                mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

					mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
            },
			wantErr: false,
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
                item: todo.TodoItem{
					Title: "", 
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
                mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)
				
					mock.ExpectRollback()
            },
			wantErr: true,
		},
		{
			name: "2nd Insert Error",
			args: args{
				listId: 1,
                item: todo.TodoItem{
					Title: "test title", 
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
                mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

					mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnError(errors.New("sone error"))
					mock.ExpectRollback()
            },
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
            testCase.mockBehavior(testCase.args, testCase.id)
			got, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
                return
            } else {
				assert.NoError(t,err)
				assert.Equal(t, testCase.id, got)
			}
			
	})
}
}