BEGIN;

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS patients (
    id          uuid NOT NULL PRIMARY KEY,
    email       varchar (320) NOT NULL UNIQUE,
    first_name  text NOT NULL,
    last_name   text NOT NULL,
    created_at  timestamptz DEFAULT NOW(),
    updated_at  timestamptz
);

CREATE TABLE IF NOT EXISTS treatment_centers (
    id          uuid NOT NULL PRIMARY KEY,
    name        text NOT NULL,
    address     text NOT NULL,
    phone       text NOT NULL,
    created_at  timestamptz DEFAULT NOW(),
    updated_at  timestamptz
);

CREATE TYPE appointment_status AS ENUM ('awaiting confirmation', 'confirmed');

CREATE TABLE IF NOT EXISTS appointments (
    id                    uuid NOT NULL PRIMARY KEY,
    treatment_center_id   uuid REFERENCES treatment_centers (id),
    start_time            timestamptz,
    created_at            timestamptz DEFAULT NOW(),
    updated_at            timestamptz
 );

CREATE TABLE IF NOT EXISTS appointment_bookings (
    id                    uuid NOT NULL PRIMARY KEY,
    appointment_id        uuid REFERENCES appointments (id),
    patient_id            uuid REFERENCES patients (id),
    status                appointment_status NOT NULL
);

CREATE UNIQUE INDEX appointment_bookings_appointment_id_patient_id_uindex
    ON appointment_bookings (appointment_id, patient_id);

COMMIT;
