{{/*
Expand the name of the chart.
*/}}
{{- define "product-store.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "product-store.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "product-store.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "product-store.commonLabels" -}}
helm.sh/chart: {{ include "product-store.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Labels
*/}}
{{- define "product-store.labels" -}}
{{ include "product-store.commonLabels" . }}
{{ include "product-store.selectorLabels" . }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "product-store.selectorLabels" -}}
app.kubernetes.io/name: {{ include "product-store.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
workbench.j4ns8i.github.com/component: product-store
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "product-store.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "product-store.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "product-store.redisHost" -}}
{{- default "redis-master" .Values.redisHost }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "product-store.redisPort" -}}
{{- default 6379 .Values.redisPort }}
{{- end }}

{{/*
Redis password Secret name
*/}}
{{- define "product-store.redisSecretName" -}}
{{- default "redis" .Values.redisPasswordName }}
{{- end }}

{{/*
Redis password Secret key
*/}}
{{- define "product-store.redisSecretPasswordKey" -}}
{{- default "redis-password" .Values.redisPasswordKey }}
{{- end }}
