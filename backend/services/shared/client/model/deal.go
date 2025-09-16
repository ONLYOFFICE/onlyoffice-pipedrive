/**
 *
 * (c) Copyright Ascensio System SIA 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package model

import "time"

type Deal struct {
	ID                int                    `json:"id"`
	Title             string                 `json:"title"`
	CreatorUserID     int                    `json:"creator_user_id"`
	OwnerID           int                    `json:"owner_id"`
	Value             float64                `json:"value"`
	PersonID          int                    `json:"person_id"`
	OrgID             int                    `json:"org_id"`
	StageID           int                    `json:"stage_id"`
	PipelineID        int                    `json:"pipeline_id"`
	Currency          string                 `json:"currency"`
	ArchiveTime       *time.Time             `json:"archive_time"`
	AddTime           time.Time              `json:"add_time"`
	UpdateTime        time.Time              `json:"update_time"`
	StageChangeTime   time.Time              `json:"stage_change_time"`
	Status            string                 `json:"status"`
	IsArchived        bool                   `json:"is_archived"`
	IsDeleted         bool                   `json:"is_deleted"`
	Probability       int                    `json:"probability"`
	LostReason        *string                `json:"lost_reason"`
	VisibleTo         int                    `json:"visible_to"`
	CloseTime         *time.Time             `json:"close_time"`
	WonTime           *time.Time             `json:"won_time"`
	LostTime          *time.Time             `json:"lost_time"`
	LocalWonDate      *string                `json:"local_won_date"`
	LocalLostDate     *string                `json:"local_lost_date"`
	LocalCloseDate    *string                `json:"local_close_date"`
	ExpectedCloseDate *string                `json:"expected_close_date"`
	LabelIDs          []int                  `json:"label_ids"`
	Origin            string                 `json:"origin"`
	OriginID          *int                   `json:"origin_id"`
	Channel           int                    `json:"channel"`
	ChannelID         string                 `json:"channel_id"`
	ACV               float64                `json:"acv"`
	ARR               float64                `json:"arr"`
	MRR               float64                `json:"mrr"`
	CustomFields      map[string]interface{} `json:"custom_fields"`
}

type DealResponse struct {
	Success bool `json:"success"`
	Data    Deal `json:"data"`
}
