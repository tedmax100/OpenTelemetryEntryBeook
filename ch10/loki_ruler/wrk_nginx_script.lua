-- script.lua
math.randomseed(os.time())

request = function()
    paths = {"/", "/666"}
    -- 隨機選擇一個路徑
    path = paths[math.random(1, #paths)]
    return wrk.format("GET", path)
end