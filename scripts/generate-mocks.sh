#!/bin/bash

cd ../api/domain
mockgen -source=entity_models.go -destination=mock_models.go -package=domain