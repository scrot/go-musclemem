package musclemem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Exercise struct {
	Owner       string  `json:"owner"`
	Workout     int     `json:"workout"`
	Index       int     `json:"index"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Repetitions int     `json:"repetitions"`
}

type exerciseService service

func (s *exerciseService) List(ctx context.Context, owner string, workout int) (*[]Exercise, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d/exercises", owner, workout)

	resp, err := s.client.send(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	xs := new([]Exercise)
	if err := dec.Decode(xs); err != nil {
		return nil, resp, err
	}

	return xs, resp, nil
}

func (c *exerciseService) Add(ctx context.Context, owner string, workout int, e Exercise) (*Exercise, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d/exercises", owner, workout)

	resp, err := c.client.send(ctx, http.MethodPost, path, e)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	respExercise := new(Exercise)
	if err := json.NewDecoder(resp.Body).Decode(respExercise); err != nil {
		return nil, resp, err
	}

	return respExercise, resp, nil
}

func (c *exerciseService) Delete(ctx context.Context, owner string, workout int, index int) (*Exercise, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d/exercises/%d", owner, workout, index)

	resp, err := c.client.send(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	respExercise := new(Exercise)
	if err := json.NewDecoder(resp.Body).Decode(respExercise); err != nil {
		return nil, nil, err
	}

	return respExercise, resp, nil
}

func (c *exerciseService) Update(ctx context.Context, owner string, workout int, index int, e Exercise) (*Exercise, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d/exercises/%d", owner, workout, index)

	resp, err := c.client.send(ctx, http.MethodPatch, path, e)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	respExercise := new(Exercise)
	if err := json.NewDecoder(resp.Body).Decode(respExercise); err != nil {
		return nil, resp, err
	}

	return respExercise, resp, nil
}

type Move string

const (
	MoveDown Move = "down"
	MoveUp   Move = "up"
	MoveSwap Move = "swap"
)

func (s *exerciseService) Move(ctx context.Context, owner string, workout int, index int, dir Move, with *int) (*http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d/exercises/%d/%s", owner, workout, index, dir)

	switch dir {
	case MoveSwap:
		if with == nil {
			return nil, errors.New("swap requires with reference")
		}

		resp, err := s.client.send(ctx, http.MethodPut, path, with)
		if err != nil {
			return resp, err
		}

		return resp, nil
	case MoveDown, MoveUp:
		if with != nil {
			return nil, errors.New(string(dir) + " requires no with reference")
		}

		resp, err := s.client.send(ctx, http.MethodPut, path, nil)
		if err != nil {
			return resp, err
		}

		return resp, nil
	default:
		return nil, errors.New("unknown move operation")
	}
}
