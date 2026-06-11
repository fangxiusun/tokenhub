package service

import (
	"sync"

	"github.com/fangxiusun/tokenhub/model"
)

var (
	channelCache     map[int]*model.Channel
	channelCacheMu   sync.RWMutex
	channelCacheInit sync.Once
)

// InitChannelCache initializes the channel cache
func InitChannelCache() {
	channels, err := model.GetAllChannels()
	if err != nil {
		return
	}

	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	channelCache = make(map[int]*model.Channel)
	for i := range channels {
		channelCache[channels[i].Id] = &channels[i]
	}
}

// GetChannelFromCache returns a channel from cache
func GetChannelFromCache(id int) *model.Channel {
	channelCacheMu.RLock()
	defer channelCacheMu.RUnlock()

	if channel, ok := channelCache[id]; ok {
		return channel
	}
	return nil
}

// UpdateChannelCache updates the channel cache
func UpdateChannelCache(ch *model.Channel) {
	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	if channelCache == nil {
		channelCache = make(map[int]*model.Channel)
	}
	channelCache[ch.Id] = ch
}

// RemoveChannelFromCache removes a channel from cache
func RemoveChannelFromCache(id int) {
	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	delete(channelCache, id)
}

// SelectChannel selects a channel for a request
func SelectChannel(group, modelName string) (*model.Channel, error) {
	abilities, err := model.GetAbilitiesByGroupAndModel(group, modelName)
	if err != nil {
		return nil, err
	}

	if len(abilities) == 0 {
		return nil, ErrNoAvailableChannel
	}

	ability := abilities[0]

	channel := GetChannelFromCache(ability.ChannelId)
	if channel == nil {
		channel, err = model.GetChannelById(ability.ChannelId)
		if err != nil {
			return nil, err
		}
		UpdateChannelCache(channel)
	}

	return channel, nil
}

// GetUserUsableGroups returns the groups a user can access
func GetUserUsableGroups(userGroup string) map[string]bool {
	groups := map[string]bool{
		"default": true,
	}
	if userGroup != "" && userGroup != "default" {
		groups[userGroup] = true
	}
	return groups
}

