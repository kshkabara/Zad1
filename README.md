# Zadanie 2 – GitHub Actions

## Schemat działania pipeline
```
Push do GitHub
↓
Checkout
↓
QEMU
↓
Buildx
↓
DockerHub Login
↓
GHCR Login
↓
Build Image
↓
Trivy Scan
↓
Build Multiarch
↓
Push do GHCR
```
---

## Budowanie obrazu wieloarchitekturowego

W celu spełnienia wymagań zadania wykorzystano Docker Buildx oraz QEMU. Dzięki temu obraz budowany jest dla dwóch architektur:

* linux/amd64
* linux/arm64

Fragment konfiguracji:

```yaml
 - name: Build and Push
        platforms: linux/amd64,linux/arm64
```

---

## Wykorzystanie cache

Dane cache przechowywane są w publicznym repozytorium DockerHub:

```text
kshka/zad2cache
```

Wykorzystano backend typu registry oraz tryb `mode=max`.

Fragment konfiguracji:

```yaml
- name: Build image for scan
    cache-from: type=registry,ref=docker.io/kshka/zad2cache:buildcache
    cache-to: type=registry,ref=docker.io/kshka/zad2cache:buildcache,mode=max
```

Takie rozwiązanie pozwala wykorzystać warstwy z poprzednich buildów i skrócić czas kolejnych procesów budowania.

---

## Test bezpieczeństwa CVE
Do skanowania podatności wykorzystano narzędzie Trivy.

Pipeline zostaje zatrzymany w przypadku wykrycia podatności:
* HIGH
* CRITICAL
Fragment konfiguracji:

```yaml
- name: Scan image with Trivy
      severity: HIGH,CRITICAL
      exit-code: 1
```

Dzięki temu obraz zostanie opublikowany wyłącznie wtedy, gdy nie zawiera podatności o wysokim lub krytycznym poziomie zagrożenia.

---

## Publikacja obrazu

Po pomyślnym zakończeniu wszystkich etapów obraz publikowany jest do GitHub Container Registry.

Do oznaczania obrazu wykorzystano tags:

```yaml
 - name: Build and Push
    tags: ghcr.io/kshkabara/zad2:latest
```
