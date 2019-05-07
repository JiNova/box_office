package main

import (
	"errors"
	"reflect"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBHandler struct {
	db *gorm.DB
}

func (handle *DBHandler) Init() {
	if handle.db != nil {
		panic("already connected")
	}

	var err error
	handle.db, err = gorm.Open("sqlite3", "/tmp/test.db")

	if err != nil {
		panic("failed connection to db")
	}
}

func (handle *DBHandler) Close() {
	_ = handle.db.Close()
}

func (handle *DBHandler) CountEntries(resources interface{}) (count int, err error) {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct:
		handle.db.Model(resources).Count(&count)
	default:
		err = errors.New("Resources must be a struct!")
	}

	return
}

func (handle *DBHandler) FillModelById(resources interface{}, id int) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct:
		if err := handle.db.First(resources, id).Error; err != nil {
			return err
		}
	default:
		return errors.New("Resources must be a struct!")
	}

	return nil
}

func (handle *DBHandler) FillModels(resources interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Slice:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}
	default:
		return errors.New("Resources must be a slice!")
	}

	return nil
}

func (handle *DBHandler) QueryModel(resources interface{}, statement string, args ...interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct, reflect.Slice:
		if err := handle.db.Where(statement, args...).Find(resources).Error; err != nil {
			return err
		}
	default:
		return errors.New("Resources must be a struct or slice!")
	}

	return nil
}

func (handle *DBHandler) QueryModelAndCount(resources interface{},
	statement string, args ...interface{}) (count int, err error) {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct, reflect.Slice:
		if err := handle.db.Where(statement, args...).Find(resources).Count(&count).Error; err != nil {
			return 0, err
		}
	default:
		return 0, errors.New("Resources must be a struct or slice!")
	}

	return
}

func (handle *DBHandler) QueryModelAndDeleteData(resources interface{}, statement string, args ...interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct, reflect.Slice:
		if err := handle.db.Where(statement, args...).Delete(resources).Error; err != nil {
			return err
		}
	}

	return nil
}

func (handle *DBHandler) CreateAssociations(resources interface{}, association string, data interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Slice:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}

		for i := 0; i < t.Len(); i++ {
			model := handle.db.Model(t.Index(i).Addr().Interface())
			field := t.Index(i).FieldByName(association).Addr().Interface()

			if err := model.Association(association).Append(field).Error; err != nil {
				panic(err)
			}
		}
	case reflect.Struct:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}

		model := handle.db.Model(t.Addr().Interface())
		//		field := t.FieldByName(association).Addr().Interface()

		if err := model.Association(association).Append(data).Error; err != nil {
			panic(err)
		}
	default:
		return errors.New("Resources must be a slice or struct!")
	}

	return nil
}

func (handle *DBHandler) LoadAssociations(resources interface{}, assocations ...string) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Slice:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}

		for _, association := range assocations {
			for i := 0; i < t.Len(); i++ {
				model := handle.db.Model(t.Index(i).Addr().Interface())
				field := t.Index(i).FieldByName(association).Addr().Interface()

				if err := model.Association(association).Find(field).Error; err != nil {
					panic(err)
				}
			}
		}
	case reflect.Struct:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}

		for _, association := range assocations {
			model := handle.db.Model(t.Addr().Interface())
			field := t.FieldByName(association).Addr().Interface()

			if err := model.Association(association).Find(field).Error; err != nil {
				panic(err)
			}
		}
	default:
		return errors.New("Resources must be a slice or struct!")
	}

	return nil
}

func (handle *DBHandler) LoadRelated(model interface{}, related interface{}) (err error) {
	t := reflect.Indirect(reflect.ValueOf(model))

	switch t.Kind() {
	case reflect.Struct, reflect.Slice:
		if err = handle.db.Model(model).Related(related).Error; err != nil {
			return
		}
	default:
		err = errors.New("Resources must be a struct or slice!")
		return
	}

	return
}
