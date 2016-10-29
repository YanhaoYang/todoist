package lib

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	SyncFailed = errors.New("Sync Failed")
)

type Sync struct {
	CollaboratorStates []interface{} `json:"collaborator_states"`
	Collaborators      []interface{} `json:"collaborators"`
	DayOrders          interface{}   `json:"day_orders"`
	DayOrdersTimestamp string        `json:"day_orders_timestamp"`
	Filters            []struct {
		Color     int    `json:"color"`
		ID        int    `json:"id"`
		IsDeleted int    `json:"is_deleted"`
		ItemOrder int    `json:"item_order"`
		Name      string `json:"name"`
		Query     string `json:"query"`
	} `json:"filters"`
	FullSync          bool   `json:"full_sync"`
	Items             []Item `json:"items"`
	Labels            Labels `json:"labels"`
	LiveNotifications []struct {
		CompletedTasks   int    `json:"completed_tasks"`
		Created          int    `json:"created"`
		DateReached      int    `json:"date_reached"`
		ID               int    `json:"id"`
		IsDeleted        int    `json:"is_deleted"`
		KarmaLevel       int    `json:"karma_level"`
		NotificationKey  string `json:"notification_key"`
		NotificationType string `json:"notification_type"`
		PromoImg         string `json:"promo_img"`
		SeqNo            int    `json:"seq_no"`
		TopProcent       int    `json:"top_procent"`
	} `json:"live_notifications"`
	LiveNotificationsLastReadID int           `json:"live_notifications_last_read_id"`
	Locations                   []interface{} `json:"locations"`
	Notes                       []struct {
		Content        string      `json:"content"`
		FileAttachment interface{} `json:"file_attachment"`
		ID             int         `json:"id"`
		IsArchived     int         `json:"is_archived"`
		IsDeleted      int         `json:"is_deleted"`
		ItemID         int         `json:"item_id"`
		Posted         string      `json:"posted"`
		PostedUID      int         `json:"posted_uid"`
		ProjectID      int         `json:"project_id"`
		UidsToNotify   interface{} `json:"uids_to_notify"`
	} `json:"notes"`
	ProjectNotes []interface{} `json:"project_notes"`
	Projects     Projects      `json:"projects"`
	Reminders    []struct {
		DateLang     string `json:"date_lang"`
		DueDateUtc   string `json:"due_date_utc"`
		ID           int    `json:"id"`
		IsDeleted    int    `json:"is_deleted"`
		ItemID       int    `json:"item_id"`
		MinuteOffset int    `json:"minute_offset"`
		NotifyUID    int    `json:"notify_uid"`
		Service      string `json:"service"`
		Type         string `json:"type"`
	} `json:"reminders"`
	SyncToken     string   `json:"sync_token"`
	TempIDMapping struct{} `json:"temp_id_mapping"`
	User          User     `json:"user"`
}

func FetchCache(token string) (Sync, error) {
	var sync Sync
	resp, err := http.PostForm(
		"https://todoist.com/API/v7/sync",
		url.Values{"token": {token}, "sync_token": {"*"}, "resource_types": {"[\"all\"]"}},
	)
	if err != nil {
		return Sync{}, SyncFailed
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Sync{}, SyncFailed
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &sync)
	if err != nil {
		return Sync{}, SyncFailed
	}
	return sync, nil
}

func LoadCache(filename string) (Sync, error) {
	var s Sync
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return s, SyncFailed
	}
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return s, SyncFailed
	}
	return s, nil
}

func SaveCache(filename string, sync Sync) error {
	buf, err := json.MarshalIndent(sync, "", "  ")
	if err != nil {
		return SyncFailed
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return SyncFailed
	}
	return nil
}
