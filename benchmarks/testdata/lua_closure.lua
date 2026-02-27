local sum = 0
local function adder(x)
  sum = sum + x
  return sum
end
for i = 0, 999 do
  adder(i)
end
return sum
