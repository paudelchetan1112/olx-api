CREATE TABLE listings (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
title TEXT NOT NULL,
description TEXT NOT NULL,
price NUMERIC(12,2) NOT NULL, 
city TEXT NOT NULL,
status TEXT NOT NULL, 
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);