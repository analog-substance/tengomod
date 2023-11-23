package os2_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/analog-substance/fileutil"
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/require"
	"github.com/analog-substance/tengomod/internal/test"
)

func TestOS2(t *testing.T) {
	rootTempDir := t.TempDir()

	lines := []interface{}{
		"line1",
		"line2",
		"line3",
	}

	tempFile1 := filepath.Join(rootTempDir, "file1.txt")
	test.Module(t, "os2").Call("write_file", tempFile1, lines).ExpectNil()

	bytes, err := os.ReadFile(tempFile1)

	require.NoError(t, err)
	require.Equal(t, "line1\nline2\nline3\n", string(bytes))

	tempFile2 := filepath.Join(rootTempDir, "file2.txt")
	data := "line1\nline2\nline3\nline4"
	test.Module(t, "os2").Call("write_file", tempFile2, data).ExpectNil()

	bytes, err = os.ReadFile(tempFile2)

	require.NoError(t, err)
	require.Equal(t, data, string(bytes))

	test.Module(t, "os2").Call("write_file", filepath.Join(rootTempDir, "nonexistent", "os2.txt"), data).ExpectTengoError()
	test.Module(t, "os2").Call("read_file_lines", tempFile2).Expect([]interface{}{"line1", "line2", "line3", "line4"})

	test.Module(t, "os2").Call("regex_replace_file", tempFile2, "line[0-9]+", "replaced").ExpectNil()

	bytes, err = os.ReadFile(tempFile2)

	require.NoError(t, err)
	require.Equal(t, "replaced\nreplaced\nreplaced\nreplaced", string(bytes))

	test.Module(t, "os2").Call("regex_replace_file", tempFile2, "line[", "replaced").ExpectTengoError()
	test.Module(t, "os2").Call("regex_replace_file", "nonexistent.txt", "line[", "replaced").ExpectTengoError()

	test.Module(t, "os2").Call("mkdir_all", tempFile2).ExpectTengoError()
	test.Module(t, "os2").Call("mkdir_all", filepath.Join(rootTempDir, "dir1", "dir2")).ExpectNil()
	require.True(t, fileutil.DirExists(filepath.Join(rootTempDir, "dir1", "dir2")))

	callRes := test.Module(t, "os2").Call("mkdir_temp", rootTempDir, "*")
	tempDir1 := callRes.Obj.(*tengo.String).Value
	require.IsType(t, "", tempDir1)
	require.True(t, fileutil.DirExists(tempDir1))

	test.Module(t, "os2").Call("copy_files", []interface{}{"nonexistent.txt"}, tempFile1).ExpectTengoError()
	test.Module(t, "os2").Call("copy_files", []interface{}{tempFile1, tempFile2}, tempDir1).ExpectNil()
	require.True(t, fileutil.FileExists(filepath.Join(tempDir1, filepath.Base(tempFile1))))
	require.True(t, fileutil.FileExists(filepath.Join(tempDir1, filepath.Base(tempFile2))))

	tempFile3 := filepath.Join(tempDir1, "file3.txt")
	test.Module(t, "os2").Call("copy_files", filepath.Join(rootTempDir, "*1.txt"), tempFile3).ExpectNil()
	require.True(t, fileutil.FileExists(tempFile3))

	test.Module(t, "os2").Call("copy_dirs", "nonexistent", filepath.Join(rootTempDir, "dir3")).ExpectTengoError()
	test.Module(t, "os2").Call("copy_dirs", []interface{}{"", ""}, filepath.Join(rootTempDir, "dir3")).ExpectTengoError()
	test.Module(t, "os2").Call("copy_dirs", tempDir1, filepath.Join(rootTempDir, "dir3")).ExpectNil()
	require.True(t, fileutil.DirExists(filepath.Join(rootTempDir, "dir3")))
}
