local count = redis.call('ZCARD', KEYS[1])

if count < 10 then
	return {}
end

return redis.call('ZPOPMIN', KEYS[1], 10)
