package service

import (
	"sync"

	"github.com/your-username/your-project/model"
)

var (
	channelCache     map[int]*model.Channel
	channelCacheMu   sync.RWMutex
	channelCacheInit sync.Once
)

// InitChannelCache initializes the channel cache
func InitChannelCache() {
	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	channels, err := model.GetAllChannels()
	if err != nil {
		return
	}

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
func UpdateChannelCache(channel *model.Channel) {
	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	if channelCache == nil {
		channelCache = make(map[int]*model.Channel)
	}
	channelCache[channel.Id] = channel
}

// RemoveChannelFromCache removes a channel from cache
func RemoveChannelFromCache(id int) {
	channelCacheMu.Lock()
	defer channelCacheMu.Unlock()

	delete(channelCache, id)
}

// SelectChannel selects a channel for a request
func SelectChannel(group, model string) (*model.Channel, error) {
	// Get abilities for the group and model
	abilities, err := model.GetAbilitiesByGroupAndModel(group, model)
	if err != nil {
		return nil, err
	}

	if len(abilities) == 0 {
		return nil, ErrNoAvailableChannel
	}

	// Select the first ability (highest priority)
	ability := abilities[0]

	// Get channel from cache or database
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
