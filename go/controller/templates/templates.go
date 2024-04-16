package templates

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	. "github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/models"
	"gorm.io/gorm"
)

func Index(c *fiber.Ctx) error {
	templates := []models.Template{}
	err := DB.Model(&models.Template{}).Preload("TemplateCustomButtons").Find(&templates).Error
	if err != nil {
		return err
	}

	return c.JSON(templates)
}

func Create(c *fiber.Ctx) error {
	request := models.Template{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Name == "" {
		return errors.New("name is required")
	}
	if request.Body == "" {
		return errors.New("body is required")
	}

	template := models.Template{
		Name:   request.Name,
		Header: request.Header,
		Body:   request.Body,
		Footer: request.Footer,
	}
	template.Validate()

	err = DB.Create(&template).Error
	if err != nil {
		return err
	}
	template.TemplateCustomButtons = []models.TemplateCustomButton{}

	for _, btn := range request.TemplateCustomButtons {
		template.CreateCustomButton(btn.Text)
	}

	return c.JSON(template)
}

func Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return errors.New("invalid id")
	}

	request := models.Template{}
	err = c.BodyParser(&request)
	if err != nil {
		return err
	}

	template := models.Template{}
	err = DB.Find(&template, id).Error
	if err != nil {
		return err
	}

	template.Name = request.Name
	template.Header = request.Header
	template.Body = request.Body
	template.Footer = request.Footer
	template.Validate()

	err = DB.Where("template_id = ?", id).Delete(&models.TemplateCustomButton{}).Error
	if err != nil {
		return err
	}

	err = DB.Save(&template).Error
	if err != nil {
		return err
	}

	template.TemplateCustomButtons = []models.TemplateCustomButton{}
	for _, btn := range request.TemplateCustomButtons {
		template.CreateCustomButton(btn.Text)
	}

	return c.JSON(template)
}

func Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return errors.New("invalid id")
	}

	template := models.Template{}
	err = DB.Find(&template, id).Error
	if err != nil {
		return err
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("template_id = ?", id).Delete(&models.TemplateCustomButton{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&template).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
