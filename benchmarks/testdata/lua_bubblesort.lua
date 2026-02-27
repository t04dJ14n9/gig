local s = {}
for i = 1, 100 do s[i] = 100 - i + 1 end
local n = #s
for i = 1, n-1 do
  for j = 1, n-1-i+1 do
    if s[j] > s[j+1] then
      s[j], s[j+1] = s[j+1], s[j]
    end
  end
end
return s[1] + s[100]
