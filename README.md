SSRF Detection and AWS Metadata Exploitation Tool

Ce programme en Go est conçu pour automatiser la détection des attaques SSRF (Server Side Request Forgery) visant les métadonnées d'instance AWS et envoyer les résultats (telles que les clés AWS et les régions) à un bot Telegram.
Fonctionnalités

    Détection de l'attaque SSRF : Le programme envoie des requêtes HTTP ciblées pour exploiter les vulnérabilités SSRF dans les services accessibles.
    Extraction des métadonnées AWS : À travers les requêtes SSRF, le programme recherche et extrait des informations sensibles telles que les clés d'accès AWS, les clés secrètes et la région de l'instance.
    Parallélisation : Utilise des goroutines et des sémaphores pour gérer plusieurs requêtes en parallèle, ce qui accélère le processus.
    Envoi des résultats à Telegram : Si des clés AWS sont trouvées, un message est envoyé à un bot Telegram avec les informations pertinentes.

Prérequis

Avant de compiler et d'exécuter le programme, assurez-vous d'avoir les éléments suivants :

    Go installé : Assurez-vous que Go est installé sur votre machine. Vous pouvez le télécharger depuis ici.
    Un bot Telegram : Vous devez avoir créé un bot Telegram et obtenir son token ainsi que l'ID du chat (voir la section ci-dessous).
    Liste d'URLs/IPs : Préparez un fichier domain.txt contenant une liste d'URLs ou d'IP à tester.

Installation

    Cloner ou télécharger le projet :

    Clonez ou téléchargez le projet sur votre machine.

git clone https://github.com/your-repository/ssrf-aws-detection.git
cd ssrf-aws-detection

Installer les dépendances :

Installez la bibliothèque fasthttp requise pour effectuer des requêtes HTTP rapides.

go get -u github.com/valyala/fasthttp

Compiler le programme :

Compilez le programme avec la commande suivante :

    go build -o ssrf_aws ssrf_aws.go

Configuration
1. Configurer le bot Telegram

Vous devez obtenir un TELEGRAM_TOKEN et un CHAT_ID pour envoyer des messages via votre bot Telegram.

    TELEGRAM_TOKEN : Obtenez ce token en créant un bot sur BotFather.
    CHAT_ID : Vous pouvez obtenir l'ID de chat en envoyant un message à votre bot et en récupérant l'ID via l'API de Telegram ou en utilisant un outil comme get_id_bot.

Dans le code, remplacez les valeurs suivantes par les vôtres :

const (
    TELEGRAM_TOKEN = "votre_telegram_token"
    CHAT_ID        = "votre_chat_id"
)

2. Préparer le fichier domain.txt

Le fichier domain.txt doit contenir une liste d'URLs ou d'adresses IP, une par ligne, que vous souhaitez tester pour détecter les vulnérabilités SSRF.

Exemple de contenu du fichier domain.txt :

example.com
192.168.1.1
subdomain.example.com

Utilisation

    Exécuter le programme :

    Une fois le programme compilé et configuré, exécutez-le avec la commande suivante :

    ./ssrf_aws

    Le programme lira les URLs/IPs depuis domain.txt, tentera d'exploiter les vulnérabilités SSRF, et envoyera un message à votre bot Telegram si des clés AWS sont découvertes.

    Sortie du programme :

    Le programme affiche des messages dans la console pour indiquer l'état des requêtes :
        [⚙️] Envoi de la requête vers les métadonnées : Indique qu'une requête SSRF est en cours.
        [🟢] Informations de métadonnées reçues : Si des clés AWS sont détectées, les informations sont affichées.
        [❌] Erreur d'accès aux métadonnées : Si une requête échoue.
        [✅] Traitement terminé : Lorsque toutes les URLs/IPs ont été traitées.

Variables de configuration

    maxWorkers : Le nombre maximal de goroutines simultanées. Par défaut, c'est 1.
    batchSize : Le nombre d'URLs/IPs traitées par lot. Par défaut, c'est 1.
    delayBetweenRequests : Le délai (en millisecondes) entre les requêtes HTTP. Par défaut, 500 ms.
    delayBetweenBatches : Le délai (en millisecondes) entre les lots de requêtes. Par défaut, 500 ms.

Vous pouvez ajuster ces variables dans le code pour optimiser les performances selon vos besoins.
Sécurité

Attention : Ce programme est conçu à des fins de test de sécurité. Assurez-vous de l'utiliser de manière éthique et légale, uniquement sur des ressources que vous avez l'autorisation de tester.
