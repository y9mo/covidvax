BEGIN;

DROP INDEX IF EXISTS appointment_bookings_appointment_id_patient_id_uindex;
DROP TABLE IF EXISTS appointment_bookings;
DROP TABLE IF EXISTS appointments;
DROP TYPE IF EXISTS appointment_status;
DROP TABLE IF EXISTS treatment_centers;
DROP TABLE IF EXISTS patients;

COMMIT;
