package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"rentacar/types"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ****************************************************************************************************
// ****************************************************************************************************
func HighlightError(err error) {

	fmt.Println("*****************************************************************************************")
	fmt.Println("*****************************************************************************************")
	fmt.Println("*****************************************************************************************")
	fmt.Println(" ")
	fmt.Println(err.Error())
	fmt.Println(" ")
	fmt.Println("*****************************************************************************************")
	fmt.Println("*****************************************************************************************")
	fmt.Println("*****************************************************************************************")

}

// ****************************************************************************************************
// ****************************************************************************************************
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ****************************************************************************************************
// concatena criterio 'where' para ser aplicado a uma query mysql
// ****************************************************************************************************
func ConcatWhere(current_where, new_criteria string) string {

	where := ""
	if current_where != "" {
		where += " and "
	}
	where += new_criteria

	return (where)
}

// ****************************************************************************************************
// salva arquivo no repositorio ftp padrao
// ****************************************************************************************************
func UploadFileToAWS_S3(localDestinationFile, ftpDestinationFile string) error {

	var err error

	imgFolder := os.Getenv("AWS_S3_IMAGES_FOLDER")

	// arquivo que foi escolhido pelo usuario
	file, err := os.Open(localDestinationFile)

	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// conexao AWS S3
	svc, err := AWS_Connect()
	if err != nil {
		return err
	}

	bucket := os.Getenv("AWS_S3_BUCKET")

	// sobe arquivo
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(imgFolder + "/" + filepath.Base(localDestinationFile)),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	if err != nil {
		return err
	}

	// verifica se chegou ok
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imgFolder + "/" + filepath.Base(localDestinationFile)),
	})
	if err != nil {
		return err
	}

	return nil
}

// ****************************************************************************************************
// copia arquivos do grupo Admin para algum grupo recem criado
// ****************************************************************************************************
func CopyAdminFilesInAWS_S3(toCopy []*types.FilesToCopyFromAdmin, processStatus *types.CloningWokgroupStatus, initialPercent int, finalPercent int) error {

	processStatus.PercentReady = strconv.Itoa(initialPercent)

	var err error

	imgFolder := os.Getenv("AWS_S3_IMAGES_FOLDER")

	// conexao AWS S3
	svc, err := AWS_Connect()
	if err != nil {
		return err
	}

	contFiles := 0
	for range toCopy {
		contFiles++
	}

	bucket := os.Getenv("AWS_S3_BUCKET")
	awsUrl := os.Getenv("AWS_S3_URL")

	for i, file := range toCopy {

		// qdo passo 1o parametro em branco (divisor= ^), frontend nao deve mostrar o 1o parametro do texto, só o 2o
		processStatus.Status = "^" + awsUrl + file.NewGroupFilename

		_, err = svc.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s/%s", bucket, imgFolder, file.AdminFilename)),
			Key:        aws.String(fmt.Sprintf("%s/%s", imgFolder, file.NewGroupFilename)),
		})
		if err != nil {
			return err
		}

		// Wait to see if the item got copied
		err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(imgFolder + "/" + file.NewGroupFilename),
		})
		if err != nil {
			return err
		}

		processStatus.PercentReady = strconv.Itoa(initialPercent + ((finalPercent - initialPercent) * i / contFiles))
		time.Sleep(300 * time.Millisecond)
	}

	return nil
}

// ****************************************************************************************************
//
//	 tenta encaixar parametros recebidos via URL (exemplo: name=Leande&city=Santos em uma Struct)
//	 exemplo tipico de uso:
//	 		URL recebida (name=Leande&city=Santos), encaixar na struct: DatatableParamsRequest
//			para verificar se foram enviados detalhes de como retornar datatable
//
// ***************************************************************************************************
func UrlToStruct(_url string, request *types.DatatableParamsRequest) error {

	var decoder = schema.NewDecoder()
	u, err := url.Parse(_url)
	if err != nil {
		HighlightError(err)
		return err
	}
	err = decoder.Decode(request, u.Query())
	if err != nil {
		HighlightError(err)
		return err
	}
	return nil
}

// *************************************************************************************************************************
// verifica se a imagem do carro escolhida é aceitavel
// *************************************************************************************************************************
func ValidateIAndUploadImageFile(c *fiber.Ctx, file_id int) error {

	var err error
	var chosenFile *multipart.FileHeader

	// FullURI exemplo ==> http://localhost:8070/workgroup/cars/43
	_url := string(c.Request().URI().FullURI())

	// verifica de qual registro se trata (car, manufacturer, etc)
	var recordType string

	if strings.Contains(_url, "/car") { // foto do carro
		recordType = "car"
	}

	if strings.Contains(_url, "/manufacturer") { // logotipo do fabricante
		recordType = "manufacturer"
	}

	if chosenFile, err = c.FormFile("chosen_image_file"); err != nil {
		return err
	}

	// exige que seja PNG
	file_type := chosenFile.Header.Get("Content-Type")
	if !strings.Contains(file_type, "image/png") {
		return errors.New("png_needed")
	}

	// grava no respositorio de imagens
	// cria nome baseado no ID do registro
	suffixWithID := fmt.Sprintf("%06d", file_id)

	// arquivo gravado localmente e arquivo que sera gravado no repositorio
	uniqueFileName := fmt.Sprintf("%v_%v_%v.png", c.Params("workgroup"), recordType, suffixWithID) // workgroup_car_000005.png, workgroup_manufacturer_000053.png, etc

	localDestinationFile := os.Getenv("LOCAL_TMP_FOLDER") + "/" + uniqueFileName // /tmp/workgroup_car_000005.png,   /tmp/workgroup_manufacturer_000053.png,   etc
	ftpDestinationFile := uniqueFileName

	// remove arquivo se previamente existir para evitar cache

	if FileExists(localDestinationFile) {
		os.Remove(localDestinationFile)
	}

	// grava localmente
	if err = c.SaveFile(chosenFile, localDestinationFile); err != nil {
		return err
	}

	// verifica tamanho
	uploadedFileAttr, _ := os.Stat(localDestinationFile)

	// max 1.5 MB
	if uploadedFileAttr.Size() > 1500000 {
		return errors.New("1.500kb_limited")
	}

	result := UploadFileToAWS_S3(localDestinationFile, ftpDestinationFile)

	return result
}

// ****************************************************************************************************
// ****************************************************************************************************
func AWS_Connect() (*s3.S3, error) {

	// sobe arquivo para S3 da AWS
	apiKey := os.Getenv("AWS_S3_APIKEY")
	secretKey := os.Getenv("AWS_S3_SECRETKEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"),
		Credentials: credentials.NewStaticCredentials(
			apiKey,
			secretKey,
			""),
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

// ****************************************************************************************************
// exclui arquivos na aws s3
// ****************************************************************************************************
func DeleteFilesInAWS_S3(toDelete []*types.FilesToDeleteInAWS) error {

	var err error

	imgFolder := os.Getenv("AWS_S3_IMAGES_FOLDER")

	// conexao AWS S3
	svc, err := AWS_Connect()
	if err != nil {
		return err
	}

	bucket := os.Getenv("AWS_S3_BUCKET")

	for _, file := range toDelete {

		fmt.Println("apagando: " + fmt.Sprintf("%s/%s", imgFolder, file.Filename))
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/%s", imgFolder, file.Filename)),
		}

		_, err = svc.DeleteObject(input)
		if err != nil {
			return err
		}

	}

	return nil
}
