package logfw

import "testing"

func TestBackupFormat(t *testing.T) {
	samples := []struct {
		name   string
		number int
		result string
	}{
		{
			name:   "test",
			number: 0,
			result: "test_0",
		},
		{
			name:   "one/test",
			number: 1,
			result: "one/test_1",
		},
		{
			name:   "one/two/test",
			number: 2,
			result: "one/two/test_2",
		},
		{
			name:   "one/two/three/test",
			number: 33,
			result: "one/two/three/test_33",
		},
		{
			name:   "program.log",
			number: 7,
			result: "program_7.log",
		},
		{
			name:   "dir1/program.log",
			number: 0,
			result: "dir1/program_0.log",
		},
		{
			name:   "dir1/dir2/program.log",
			number: 125,
			result: "dir1/dir2/program_125.log",
		},
	}
	var bf backupFormat
	for _, sample := range samples {
		bf.SetFileName(sample.name)
		backupName := bf.BackupName(sample.number)
		if backupName != sample.result {
			t.Fatalf("%s != %s", backupName, sample.result)
		}
	}
}

func TestBackupFormatNumberLen(t *testing.T) {
	samples := []struct {
		name   string
		number int
		numlen int
		result string
	}{
		{
			name:   "test",
			number: 1,
			numlen: 1,
			result: "test_1",
		},
		{
			name:   "test",
			number: 1,
			numlen: 2,
			result: "test_01",
		},
		{
			name:   "test",
			number: 1,
			numlen: 3,
			result: "test_001",
		},
		{
			name:   "test",
			number: 1,
			numlen: 4,
			result: "test_0001",
		},
		{
			name:   "test.log",
			number: 1,
			numlen: 1,
			result: "test_1.log",
		},
		{
			name:   "test.log",
			number: 1,
			numlen: 2,
			result: "test_01.log",
		},
		{
			name:   "test.log",
			number: 1,
			numlen: 3,
			result: "test_001.log",
		},
		{
			name:   "test.log",
			number: 1,
			numlen: 4,
			result: "test_0001.log",
		},
		{
			name:   "master.log",
			number: 27,
			numlen: 4,
			result: "master_0027.log",
		},
		{
			name:   "home/test.log",
			number: 24,
			numlen: 4,
			result: "home/test_0024.log",
		},
		{
			name:   "home/master.log",
			number: 123,
			numlen: 4,
			result: "home/master_0123.log",
		},
	}
	var bf backupFormat
	for _, sample := range samples {
		bf.SetFileName(sample.name)
		bf.SetNumberLen(sample.numlen)
		backupName := bf.BackupName(sample.number)
		if backupName != sample.result {
			t.Fatalf("%s != %s", backupName, sample.result)
		}
	}
}

func TestCountDigits(t *testing.T) {
	check := func(x, result int) {
		if cd := countDigits(x); cd != result {
			t.Fatalf("%d != %d", cd, result)
		}
	}
	for i := 0; i < 10; i += 1 {
		check(i, 1)
	}
	for i := 10; i < 100; i += 10 {
		check(i, 2)
	}
	for i := 100; i < 1000; i += 100 {
		check(i, 3)
	}
	for i := 1000; i < 10000; i += 1000 {
		check(i, 4)
	}
	for i := 10000; i < 100000; i += 10000 {
		check(i, 5)
	}
}
