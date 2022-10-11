package repository

import (
	"database/sql"
	"errors"
	"fmt"
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

var randomDate = time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)

func TestAddCity(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		city             string
		subscriptionDate time.Time
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO subscription").
					WithArgs("Moscow", randomDate).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			input: args{
				city:             "Moscow",
				subscriptionDate: randomDate,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.AddCity(tt.input.city, tt.input.subscriptionDate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fmt.Println(err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAvgTempByCityId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		id int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    float64
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"weather"}).
					AddRow(15.0)

				mock.ExpectQuery("SELECT AVG from weather WHERE (.+)").WithArgs("2").WillReturnRows(rows)

			},
			input: args{
				id: 2,
			},
			want: 15.0,
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT AVG(weather) from weather WHERE (.+)").WithArgs("city2").WillReturnRows(rows)

			},
			input: args{
				id: 404,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAvgTempByCityId(tt.input.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetSubscriptionList(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []weather.Subscription
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "city", "subscription_date"}).
					AddRow(1, "city1", randomDate).
					AddRow(2, "city2", randomDate).
					AddRow(3, "city3", randomDate)

				mock.ExpectQuery("SELECT (.+) FROM subscription").WillReturnRows(rows)
			},
			want: []weather.Subscription{
				{1, "city1", randomDate},
				{2, "city2", randomDate},
				{3, "city3", randomDate},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetSubscriptionList()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMoveOldDataToArchive(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		dateForDelete time.Time
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("WITH moved_rows AS (DELETE FROM weather WHERE (.+) RETURNING (.+)) INSERT INTO weather_archive SELECT * FROM moved_rows").
					WithArgs(randomDate).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			input: args{
				dateForDelete: randomDate,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.MoveOldDataToArchive(tt.input.dateForDelete)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCityId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		city string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(2)

				mock.ExpectQuery("SELECT id from subscription WHERE (.+)").WithArgs("city2").WillReturnRows(rows)

			},
			input: args{
				city: "city2",
			},
			want: 2,
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT id from subscription WHERE (.+)").WithArgs("city2").WillReturnRows(rows)

			},
			input: args{
				city: "city2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetCityId(tt.input.city)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddWeatherByCityId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		id          int
		date        time.Time
		temperature int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO weather").
					WithArgs(1, randomDate, 18).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			input: args{
				id:          1,
				date:        randomDate,
				temperature: 18,
			},
		},
		{
			name: "Not found",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO weather").
					WithArgs(404, randomDate, 404).WillReturnError(errors.New("not found"))

				mock.ExpectRollback()
			},
			input: args{
				id:          404,
				date:        randomDate,
				temperature: 404,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.AddWeatherByCityId(tt.input.id, tt.input.date, tt.input.temperature)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fmt.Println(err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteCityById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewApiPostgres(db)

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM weather WHERE").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec("DELETE FROM subscription WHERE").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			input: args{
				id: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM weather WHERE").
					WithArgs(404).WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			input: args{
				id: 404,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteCityById(tt.input.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
