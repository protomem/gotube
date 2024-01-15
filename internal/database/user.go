package database

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/model"
)

func (db *DB) GetUser(ctx context.Context, id model.ID) (model.User, error) {
	const op = "database.GetUser"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	user, err := db.getUserByField(ctx, Field{Name: "id", Value: id})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (db *DB) GetUserByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "database.GetUserByNickname"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	user, err := db.getUserByField(ctx, Field{Name: "nickname", Value: nickname})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "database.GetUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	user, err := db.getUserByField(ctx, Field{Name: "email", Value: email})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

type InsertUserDTO struct {
	Nickname string
	Password string
	Email    string
}

func (db *DB) InsertUser(ctx context.Context, dto InsertUserDTO) (model.ID, error) {
	const op = "database.InsertUser"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO users(nickname, password, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	args := []any{dto.Nickname, dto.Password, dto.Email}

	var id model.ID

	if err := db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrUserAlreadyExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

type UpdateUserDTO struct {
	Nickname    *string
	Password    *string
	Email       *string
	AvatarPath  *string
	Description *string
}

func (db *DB) UpdateUser(ctx context.Context, id model.ID, dto UpdateUserDTO) error {
	const op = "database.UpdateUserByNickname"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	fields := make([]Field, 0, 5)

	if dto.Nickname != nil {
		fields = append(fields, Field{Name: "nickname", Value: *dto.Nickname})
	}
	if dto.Password != nil {
		fields = append(fields, Field{Name: "password", Value: *dto.Password})
	}
	if dto.Email != nil {
		fields = append(fields, Field{Name: "email", Value: *dto.Email})
	}
	if dto.AvatarPath != nil {
		fields = append(fields, Field{Name: "avatar_path", Value: *dto.AvatarPath})
	}
	if dto.Description != nil {
		fields = append(fields, Field{Name: "description", Value: *dto.Description})
	}

	if err := db.updateUserByField(ctx, Field{Name: "id", Value: id}, fields); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DB) DeleteUserByNickname(ctx context.Context, nickname string) error {
	const op = "database.DeleteUserByNickname"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	if err := db.deleteUserByField(ctx, Field{Name: "nickname", Value: nickname}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DB) getUserByField(ctx context.Context, field Field) (model.User, error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE %s = $1 LIMIT 1`, field.Name)
	args := []any{field.Value}

	var user model.User

	if err := db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		if IsNoRows(err) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (db *DB) updateUserByField(ctx context.Context, byFiled Field, fields []Field) error {
	counter := 1
	query := `UPDATE users SET updated_at = now()`
	args := []any{byFiled.Value}

	for _, f := range fields {
		counter++
		query += fmt.Sprintf(`, %s = $%d`, f.Name, counter)
		args = append(args, f.Value)
	}

	query += fmt.Sprintf(` WHERE %s = $%d`, byFiled.Name, 1)

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) deleteUserByField(ctx context.Context, field Field) error {
	query := fmt.Sprintf(`DELETE FROM users WHERE %s = $1`, field.Name)
	args := []any{field.Value}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
