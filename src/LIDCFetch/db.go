package main

import (
	"fmt"

	"github.com/antonholmquist/jason"
	"github.com/mxk/go-sqlite/sqlite3"
)

func Save(db *sqlite3.Conn, table string, uid string, obj *jason.Object) error {
	handle_err := func(col string, e error) {
		if e != nil {
			logger.Error("Error inserting %v into %v: %v", col, table, e.Error())
		}
	}

	err := db.Exec(fmt.Sprintf("insert or ignore into %v ( uid ) values (?)", table), uid)
	handle_err("", err)
	for k, v := range obj.Map() {
		if k == "uid" {
			continue
		}
		// Are we a string or a number?
		if i, err := v.Int64(); err == nil {
			err = db.Exec(fmt.Sprintf("update %v set %v = ? where uid = ?", table, k), i, uid)
			handle_err(k, err)
		}
		if f, err := v.Float64(); err == nil {
			err = db.Exec(fmt.Sprintf("update %v set %v = ? where uid = ?", table, k), f, uid)
			handle_err(k, err)
		}
		if s, err := v.String(); err == nil {
			err = db.Exec(fmt.Sprintf("update %v set %v = ? where uid = ?", table, k), s, uid)
			handle_err(k, err)
		}
		if arr, err := obj.GetFloat64Array(k); err == nil {
			s := fmt.Sprintf("%v", arr)
			err = db.Exec(fmt.Sprintf("update %v set %v = ? where uid = ?", table, k), s, uid)
			handle_err(k, err)
		}
	}
	return err
}

var create_series_table = `
create table if not exists series (
  uid text primary key,
  series_instance_uid text,
  study_instance_uid text,
  patient_name text,
  patient_id text,
  manufacturer text,
  manufacturer_model_name text,
  patient_sex text,
  patient_age text,
  ethnic_group text,
  contrast_bolus_agent text,
  body_part_examined text,
  scan_options text,
  slice_thickness float,
  kvp float,
  data_collection_diameter float,
  software_versions text,
  reconstruction_diameter float,
  gantry_detector_tilt float,
  table_height float,
  rotation_direction text,
  exposure_time float,
  xray_tube_current float,
  exposure float,
  convolution_kernel text,
  patient_position text,
  image_position_patient text,
  image_orientation_patient text,
  filename text )
`
var create_read_table = `
create table if not exists reads (
  uid text primary key,
  nodule_uid text,
  normalized_nodule_id int,
  id text,
  centroid text,
  centroidLPS text,
  point_count int,
  label_value int,
  filled int,
  subtlety int,
  internalStructure int,
  calcification int,
  sphericity int,
  margin int,
  lobulation int,
  spiculation int,
  texture int,
  malignancy int  
)
`

var create_nodule_table = `
create table if not exists nodules (
  uid text primary key,
  series_uid text,
  normalized_nodule_id int
)
`

var create_measure_table = `
create table if not exists measures (
  uid text primary key,
  nodule_uid text,
  read_uid text,
  command_line text,
  false_negative_error float,
  dice_coefficient float,
  volume_similarity float,
  false_positive_error float,
  mean_overlap float,
  union_overlap float,
  jaccard_coefficient float,
  hausdorff_distance float,
  average_hausdorff_distance float
)`

var create_tables = map[string]string{"series": create_series_table, "reads": create_read_table, "nodules": create_nodule_table, "measures": create_measure_table}

var db_migrations = []string{}
