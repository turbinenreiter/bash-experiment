a()
{
	x=$((x - 1))
	return 5
}

b()
{
	x=$((x - 1))
	a
	echo a returns $?
	return 4
}

c()
{
	x=$((x - 1))
	b
	echo b returns $?
	return 3
}

d()
{
	x=$((x - 1))
	c
	echo c returns $?
	return 2
}

e()
{
	d
	echo d returns $?
	echo in e
	x=$((x - 1))
	return $x
}

f()
{
	e
	echo e returned $?
	echo x is $x
	return 0
}

x=30
f

# make sure unsetting a local variable preserves the `local' attribute
f1()
{
	local zz
	zz=abcde
	echo $zz
	unset zz
	zz=defghi
	echo $zz
}

zz=ZZ
echo $zz
f1
echo $zz

unset -f f1
f1()
{
	return 5
}

( f1 )
echo $?

