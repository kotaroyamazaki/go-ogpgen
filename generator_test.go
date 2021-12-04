package ogpgen

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func init() {
	img := image.NewNRGBA(image.Rect(0, 0, 1200, 630))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{255, 255, 255, 255}), image.Point{}, draw.Src)
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err = png.Encode(f, img); err != nil {
		panic(err)
	}

}

func Test_ogpgen_AttachImage(t *testing.T) {
	type fields struct {
		img     *image.RGBA
		quality int
		height  int
		width   int
	}
	type args struct {
		params *ImageCompositionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				img:     image.NewRGBA(image.Rect(0, 0, 1200, 630)),
				quality: 100,
				height:  1200,
				width:   630,
			},
			args: args{
				params: &ImageCompositionParams{
					ImagePath: "test.png",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to open image when image path is empty",
			fields: fields{
				img:     image.NewRGBA(image.Rect(0, 0, 1200, 630)),
				quality: 100,
				height:  100,
				width:   100,
			},
			args: args{
				params: &ImageCompositionParams{
					ImagePath: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ogpgen{
				img:     tt.fields.img,
				quality: tt.fields.quality,
				height:  tt.fields.height,
				width:   tt.fields.width,
			}
			if err := c.AttachImage(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("ogpgen.AttachImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ogpgen_AttachText(t *testing.T) {
	type fields struct {
		img     *image.RGBA
		quality int
		height  int
		width   int
	}
	type args struct {
		params *TextCompositionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				img:     image.NewRGBA(image.Rect(0, 0, 1200, 630)),
				quality: 100,
				height:  1200,
				width:   630,
			},
			args: args{
				params: &TextCompositionParams{
					Text: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to attach text when text is empty",
			fields: fields{
				img:     image.NewRGBA(image.Rect(0, 0, 1200, 630)),
				quality: 100,
				height:  1200,
				width:   630,
			},
			args: args{
				params: &TextCompositionParams{
					Text: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ogpgen{
				img:     tt.fields.img,
				quality: tt.fields.quality,
				height:  tt.fields.height,
				width:   tt.fields.width,
			}
			if err := c.AttachText(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("ogpgen.AttachText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
