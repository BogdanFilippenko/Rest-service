package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service/internal/model"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// работа с базой
type Repository struct {
	db *sql.DB
}

// создаем новый репозиторий
func New(db *sql.DB) *Repository {

	return &Repository{db: db}
}

func parseDate(monthYear string) (time.Time, error) {
	return time.Parse("01-2006", monthYear)
}
//CREATE
func (r *Repository) Create(ctx context.Context, sub model.Subscription) (int, error) {
	start, err := parseDate(sub.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date: %w", err)
	}

	var end sql.NullTime
	if sub.EndDate != nil {
		parsedEnd, err := parseDate(*sub.EndDate)
		if err != nil {
			return 0, fmt.Errorf("invalid end_date: %w", err)
		}
		end = sql.NullTime{Time: parsedEnd, Valid: true}
	}

	query :=
		`INSERT INTO subscription (service_name, price, user_id, start_date, end_date )
 VALUES ($1,$2,$3,$4,$5)
 RETURNING id
 `
	var id int
	err = r.db.QueryRowContext(ctx, query, sub.ServiceName, sub.Price, sub.UserId, start, end).Scan(&id)
	if err != nil{
		return 0, err
	} 
	return id, nil
}
//FINDBYID
func (r *Repository) FindByID(ctx context.Context, id int)(model.Subscription, error){
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE id = $1 `

	var sub model.Subscription
	var start time.Time
	var end sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserId, &start, &end)
	if err!= nil{
		if errors.Is(err, sql.ErrNoRows) {
			return sub, errors.New("not found")
		}
		return sub, err

	}

	sub.StartDate = start.Format("01-2006")
	if end.Valid{
		endStr := end.Time.Format("01-2006")
		sub.EndDate = &endStr
	}
	return sub, nil

}
//LIST
func(r *Repository)List(ctx context.Context)([]model.Subscription, error){
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription`

	rows, err:=r.db.QueryContext(ctx, query)
	if err!= nil{
		return nil, err 
	}

	defer rows.Close()

	var subs []model.Subscription
	for rows.Next(){
		var sub model.Subscription
	var start time.Time
	var end sql.NullTime

	if err:= rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserId, &start, &end); err !=nil{
		return nil, err
	}
	sub.StartDate = start.Format("01-2006")
		if end.Valid {
			endStr := end.Time.Format("01-2006")
			sub.EndDate = &endStr
		}
		subs = append(subs, sub)

	}
	if err := rows.Err(); err!= nil{
		return nil, err
	}
	return subs, nil
}

//UPDATE
func (r *Repository) Update(ctx context.Context, id int, sub model.Subscription) error {
	start, err := parseDate(sub.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}

	var end sql.NullTime
	if sub.EndDate != nil {
		parsedEnd, err := parseDate(*sub.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date: %w", err)
		}
		end = sql.NullTime{Time: parsedEnd, Valid: true}
	}

	query := `
		UPDATE subscription 
		SET service_name = $1, price = $2,user_id = $3, start_date = $4, end_date = $5 
		WHERE id = $6`

	res, err := r.db.ExecContext(ctx, query, sub.ServiceName, sub.Price, sub.UserId, start, end, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}
//DELETE
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscription WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err!=nil{
	return err
	}
	rowsAffected, err := res.RowsAffected()
	if err!=nil{
		return err
	}
	if rowsAffected == 0{
		return errors.New("not found")
	}
	return nil
}

func (r *Repository) GetForCostCalculation(ctx context.Context, UserID uuid.UUID, ServiceName string) ([]model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
	          FROM subscription WHERE user_id = $1 AND service_name = $2`

	rows, err := r.db.QueryContext(ctx, query, UserID, ServiceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		var start time.Time
		var end sql.NullTime

		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserId, &start, &end); err != nil {
			return nil, err
		}
		sub.StartDate = start.Format("01-2006")
		if end.Valid {
			endStr := end.Time.Format("01-2006")
			sub.EndDate = &endStr
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
