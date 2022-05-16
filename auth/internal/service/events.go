package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"sort"
	"strconv"
	"strings"

	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/jmoiron/sqlx"
)

type EventsService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (e *EventsService) GetAllEvents(ctx context.Context) (models.EventsResponse, error) {
	tags, err := e.repo.Tags.AllTags(ctx)
	if err != nil {
		return models.EventsResponse{}, err
	}

	res, err := e.repo.Events.AllEvents(ctx)
	if err != nil {
		return models.EventsResponse{}, err
	}
	for i, r := range res {
		var t string
		at := strings.Split(r.Tags, ";")
		if len(at) == 0 {
			t = ""
			res[i].Tags = t
		} else {
			aTags, err := e.repo.Tags.AllUnfavoriteTagsTags(ctx, at)
			if err != nil {
				return models.EventsResponse{}, err
			}
			t = ""
			for _, a := range aTags {
				if t == "" {
					t = t + a.Title
				} else {
					t = t + ", " + a.Title
				}
			}
			t = t + "."
			res[i].Tags = t
		}
	}

	result := models.EventsResponse{
		AllEvents: len(res),
		AllTags:   len(tags),
		Events:    res,
	}
	return result, nil
}

func (e *EventsService) GetEvents(ctx context.Context, req []models.AllTags, userID string) ([]models.MyEvents, error) {
	var (
		arF  []int64
		arUF []int64
	)
	//UserTags
	usrInfo, err := e.repo.Auth.UserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	if usrInfo.FTags != "" {
		ft := strings.Split(usrInfo.FTags, ";")
		arF = make([]int64, len(ft))

		for j := range arF {
			v, err := strconv.ParseInt(ft[j], 10, 64)
			if err != nil {
				return nil, err
			}
			arF[j] = v
		}
	}

	if usrInfo.UFTags != "" {
		uft := strings.Split(usrInfo.UFTags, ";")
		arUF = make([]int64, len(uft))

		for j := range arUF {
			v, err := strconv.ParseInt(uft[j], 10, 64)
			if err != nil {
				return nil, err
			}
			arUF[j] = v
		}
	}

	//AllTags
	tags, err := e.repo.Tags.AllTags(ctx)
	if err != nil {
		return nil, err
	}
	tmap := make(map[int64]int)
	for _, t := range tags {
		tmap[t.ID] = t.Val
	}

	//AllEvents
	res, err := e.repo.Events.AllEvents(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.MyEvents, 0)
	for _, r := range res {
		var (
			myE models.MyEvents
			ary []int64
			at  []string
			t   string
		)
		if r.Tags != "" {
			at = strings.Split(r.Tags, ";")
			ary = make([]int64, len(at))

			for j := range ary {
				v, err := strconv.ParseInt(at[j], 10, 64)
				if err != nil {
					return nil, err
				}
				ary[j] = v
			}
		}

		myE.Title = r.Title
		myE.Address = r.Address
		myE.PNG = r.PNG
		myE.Val = 0

		for _, a := range ary {
			for _, mr := range req {
				if a == mr.ID {
					myE.Val += tmap[a]
				}
			}

			for _, auf := range arUF {
				if a == auf {
					myE.Val -= tmap[a] / 2
				}
			}

			for _, af := range arF {
				if a == af {
					myE.Val += tmap[a] / 2
				}
			}
		}

		if len(at) == 0 {
			t = ""
			myE.Tags = t
		} else {
			aTags, err := e.repo.Tags.AllUnfavoriteTagsTags(ctx, at)
			if err != nil {
				return nil, err
			}
			t = ""
			for _, a := range aTags {
				if t == "" {
					t = t + a.Title
				} else {
					t = t + ", " + a.Title
				}
			}
			t = t + "."
		}
		myE.Tags = t

		result = append(result, myE)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Val > result[j].Val
	})

	return result, nil
}

func updateRedis(ctx context.Context, req []models.AllTags, userID string) error {

	return nil
}

func NewEventsService(repo *repository.Repository, conn *sqlx.DB) *EventsService {
	return &EventsService{
		repo: repo,
		conn: conn,
	}
}
