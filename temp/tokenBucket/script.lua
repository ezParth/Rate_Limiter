local key = KEYS[1]
local now = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local capacity = tonumber(ARGV[3])

-- get bucket state
local data = redis.call("HMGET", key, "tokens", "timestamp")
local tokens = tonumber(data[1])
local last_ts = tonumber(data[2])

if tokens == nil then
    tokens = capacity
    last_ts = now
end

-- refill tokens
local delta = math.max(0, now - last_ts)
tokens = math.min(capacity, tokens + delta * refill_rate)

-- check request
if tokens < 1 then
    redis.call("HMSET", key, "tokens", tokens, "timestamp", now)
    redis.call("EXPIRE", key, 60)
    return 0 -- reject
else
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "timestamp", now)
    redis.call("EXPIRE", key, 60)
    return 1 -- accept
end
