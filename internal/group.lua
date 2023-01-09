local count = redis.call('ZCARD', KEYS[1])

if count < 10 then
	return 0
end

local room = redis.call('SISMEMBER', KEYS[2], KEYS[3])

if room ~= 0 then
    return 0
end

redis.call('SADD', KEYS[2], KEYS[3])

local p = redis.call('ZPOPMIN', KEYS[1], 10)

for _, user in ipairs(p) do
    redis.call('SADD', KEYS[3], user)
end

return 1
