package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES ($1, $2, $3, $4)`

	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении посылки: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	query := "SELECT * FROM parcel WHERE number = $1"

	err := s.db.QueryRow(query, number).Scan(
		&p.Number,
		&p.Client,
		&p.Status,
		&p.Address,
		&p.CreatedAt,
	)

	if err != nil {
		return p, fmt.Errorf("ошибка получения посылки: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	query := "SELECT * FROM parcel WHERE client = $1"

	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса посылок клиента: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Parcel
		if err := rows.Scan(
			&p.Number,
			&p.Client,
			&p.Status,
			&p.Address,
			&p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := "UPDATE parcel SET status = $1 WHERE number = $2"

	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса: %w", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	// Проверяем статус посылки
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("ошибка получения посылки для проверки статуса: %w", err)
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("изменение адреса возможно только для посылок со статусом 'registered'")
	}

	query := "UPDATE parcel SET address = $1 WHERE number = $2"
	_, err = s.db.Exec(query, address, number)
	if err != nil {
		return fmt.Errorf("ошибка обновления адреса: %w", err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	// Проверяем статус посылки
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("ошибка получения посылки для проверки статуса: %w", err)
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("удаление возможно только для посылок со статусом 'registered'")
	}

	query := "DELETE FROM parcel WHERE number = $1"
	_, err = s.db.Exec(query, number)
	if err != nil {
		return fmt.Errorf("ошибка удаления посылки: %w", err)
	}
	return nil
}
