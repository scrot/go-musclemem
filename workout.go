package musclemem

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Workout struct {
	Owner string `json:"owner"`
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type workoutService service

func (s *workoutService) List(ctx context.Context, owner string) (*[]Workout, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts", owner)

	resp, err := s.client.send(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	ws := new([]Workout)
	if err := dec.Decode(ws); err != nil {
		return nil, resp, err
	}

	return ws, resp, nil
}

func (s *workoutService) Add(ctx context.Context, owner string, w Workout) (*Workout, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts", owner)

	resp, err := s.client.send(ctx, http.MethodPost, path, w)
	if err != nil {
		return nil, resp, err
	}

	respWorkout := new(Workout)
	if err := json.NewDecoder(resp.Body).Decode(respWorkout); err != nil {
		return nil, resp, err
	}

	return respWorkout, resp, nil
}

func (s *workoutService) Delete(ctx context.Context, owner string, index int) (*Workout, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d", owner, index)

	resp, err := s.client.send(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	respWorkout := new(Workout)
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(respWorkout); err != nil {
		return nil, nil, err
	}

	return respWorkout, resp, nil
}

func (s *workoutService) Update(ctx context.Context, owner string, index int, w Workout) (*Workout, *http.Response, error) {
	path := fmt.Sprintf("/users/%s/workouts/%d", owner, index)

	resp, err := s.client.send(ctx, http.MethodPatch, path, w)
	if err != nil {
		return nil, nil, err
	}

	respWorkout := new(Workout)
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(respWorkout); err != nil {
		return nil, nil, err
	}

	return respWorkout, resp, nil
}
