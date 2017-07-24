package recordio_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ondrej-smola/mesos-go-http/lib/codec/framing/recordio"
)

func TestRecordio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recordio Suite")
}

var _ = Describe("RecordIOReader", func() {

	It("Test", func() {
		for i, tt := range []struct {
			in     string
			out    []byte
			fwd, n int
			err    error
		}{
			{"1\na0\n1\nb", []byte("a"), 0, 1, nil},
			{"1\na0\n1\nb", []byte("b"), 1, 1, nil},
			{"1\na", []byte{}, 0, 0, nil},
			{"2\nab", []byte("a"), 0, 1, nil},
			{"2\nab", []byte("ab"), 0, 2, nil},
			{"2\nab", []byte("b"), 1, 1, nil},
			{"2\nab", []byte{'b', 0}, 1, 1, nil},
			{"2\nab", []byte{0}, 2, 0, io.EOF},
			{ // size = (2 << 63) + 1
				"18446744073709551616\n", []byte{0}, 0, 0, &strconv.NumError{
					Func: "ParseUint",
					Num:  "18446744073709551616",
					Err:  strconv.ErrRange,
				},
			},
		} {
			type expect struct {
				p   []byte
				n   int
				err error
			}
			r := recordio.New(strings.NewReader(tt.in))
			if n, err := r.Read(make([]byte, tt.fwd)); err != nil || n != tt.fwd {
				Fail(fmt.Sprintf("test #%d: failed to read forward %d bytes: %v", i, n, err))
			}

			want := expect{[]byte(tt.out), tt.n, tt.err}
			got := expect{p: make([]byte, len(tt.out))}
			if got.n, got.err = r.Read(got.p); !reflect.DeepEqual(got, want) {
				Fail(fmt.Sprintf("test #%d: got: %+v, want: %+v", i, got, want))
			}

		}
	})

	Measure("Benchmark", func(b Benchmarker) {
		count := 100000
		var buf bytes.Buffer
		genRecords(&buf)
		r := recordio.New(&buf)
		p := make([]byte, 256)

		total := 0
		runtime := b.Time("runtime", func() {
			for i := 0; i < count; i++ {
				if n, err := r.Read(p); err != nil && err != io.EOF {
					Fail(fmt.Sprint(err))
				} else {
					total += n
				}
			}
		})

		b.RecordValue("MB/s", float64(total)/runtime.Seconds()/1000000)
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.02))
	}, 5)
})

func genRecords(w io.Writer) {
	rnd := rng{rand.New(rand.NewSource(0xdeadbeef))}
	buf := make([]byte, 2<<12)
	for i := 0; i < cap(buf); i++ {
		sz := rnd.Intn(cap(buf))
		n, err := rnd.Read(buf[:sz])
		if err != nil {
			Fail(fmt.Sprint(err))
		}
		header := strconv.FormatInt(int64(n), 10) + "\n"
		if _, err = io.WriteString(w, header); err != nil {
			Fail(fmt.Sprint(err))
		} else if _, err = w.Write(buf[:n]); err != nil {
			Fail(fmt.Sprint(err))
		}
	}
}

type rng struct{ *rand.Rand }

func (r rng) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i += 7 {
		val := r.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
	return len(p), nil
}
