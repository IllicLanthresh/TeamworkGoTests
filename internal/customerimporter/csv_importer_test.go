package customerimporter

import (
	"io/ioutil"
	"log"
	"testing"
)

func Test_NewCsvCustomerImporter(t *testing.T) {
	type args struct {
		csvPath  string
		emailKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid csv path",
			args: struct {
				csvPath  string
				emailKey string
			}{
				csvPath:  "foo",
				emailKey: "",
			},
			wantErr: true,
		},
		{
			name: "path does not exist",
			args: struct {
				csvPath  string
				emailKey string
			}{
				csvPath:  "foo.csv",
				emailKey: "",
			},
			wantErr: true,
		},
		{
			name: "no email key",
			args: struct {
				csvPath  string
				emailKey string
			}{
				csvPath:  "../../test/data/importer/foo.csv",
				emailKey: "",
			},
			wantErr: true,
		},
		{
			name: "all good",
			args: struct {
				csvPath  string
				emailKey string
			}{
				csvPath:  "../../test/data/importer/foo.csv",
				emailKey: "foo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCsvCustomerImporter(tt.args.csvPath, tt.args.emailKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCsvCustomerImporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_csvCustomerImporter_emailAddressesGenerator(t *testing.T) {
	type fields struct {
		csvPath  string
		emailKey string
	}
	tests := []struct {
		name               string
		fields             fields
		wantEmailAddresses bool
		wantErr            bool
	}{
		{
			name: "no headers",
			fields: fields{
				csvPath:  "../../test/data/importer/foo.csv",
				emailKey: "foo",
			},
			wantEmailAddresses: false,
			wantErr:            true,
		},
		{
			name: "header not found",
			fields: fields{
				csvPath:  "../../test/data/importer/foo-headers.csv",
				emailKey: "quux",
			},
			wantEmailAddresses: false,
			wantErr:            true,
		},
		{
			name: "all good",
			fields: fields{
				csvPath:  "../../test/data/importer/foo-headers.csv",
				emailKey: "bar",
			},
			wantEmailAddresses: true,
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp := &csvCustomerImporter{
				csvPath:  tt.fields.csvPath,
				emailKey: tt.fields.emailKey,
			}
			gotEmailAddresses, err := imp.emailAddressesGenerator()
			if (err != nil) != tt.wantErr {
				t.Errorf("emailAddressesGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotEmailAddresses != nil) != tt.wantEmailAddresses {
				t.Errorf("emailAddressesGenerator() gotEmailAddresses = %v, want %v", gotEmailAddresses, tt.wantEmailAddresses)
			}
		})
	}
}

func Benchmark_csvCustomerImporter_CustomerCountByDomain(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	type args struct {
		csvPath  string
		emailKey string
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			name: "sorting 3k",
			args: args{
				csvPath:  "../../test/data/importer/customers.csv",
				emailKey: "email",
			},
		},
		{
			name: "sorting 1m",
			args: args{
				csvPath:  "../../test/data/importer/emails-1m.csv",
				emailKey: "email",
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			imp, err := NewCsvCustomerImporter(bm.args.csvPath, bm.args.emailKey)
			if err != nil {
				b.Error(err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := imp.CustomerCountByDomain()
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}
