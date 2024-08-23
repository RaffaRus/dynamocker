# <p5 style="color:#213F85">Dyna</p5><p5 style="color:#C2AA7E">Mocker</p5>
DynaMocker is a simple mocker of a web server which allows you to dynamically modify the `*.json` file (containing the mock response that you want to fetch to your client) in real time. For those who are found of mice, the app comes with a simple UI, that allows you to add, modify and delete the behavior of each mock API response.

The app can be used through *docker-compose* or through *helm* chart.

## <img src="https://img.icons8.com/?size=100&id=Wln8Z3PcXanx&format=png&color=000000" alt="MarineGEO circle logo" style="height: 20px; width:20px;"/> Docker-compose
The project is made of two images: one for the back-end and another one for the front-end. You have to choose:
- the folder in which your mockApis will be stored : <full_path_to_you_folder>
- the port on which the backend will be hosted: <your_be_port>
- the port on which the frontend will be reachable: <your_fe_port>

Create the `.env` file in which you insert your chosen env variables:
```
MOCK_API_FOLDER=<full_path_to_you_folder>
FE_PORT=<your_be_port>
BE_PORT=<your_fe_port>
```
Create the following `docker-compose.yml` file:
``` yml
services:
  dynamocker-be:
    image: raffarus/dynamocker-be:0.0.1
    ports:
      - ${BE_PORT}:8150
    volumes:
      - ${MOCK_API_FOLDER}:/mocks
    environment:
      - BE_PORT=${BE_PORT}
  dynamocker-fe:
    image: raffarus/dynamocker-fe:0.0.1
    ports:
      - ${FE_PORT}:8151
    environment:
      - FE_PORT=${FE_PORT}
      - BE_PORT=${BE_PORT}
    depends_on:
      - dynamocker-be
```
Then use docker-compose to launch the containers:
```
$ docker compose -f ./docker-compose.yml --env-file ./.env up -d
```
You should be able to access the UI to http://localhost:{FE_PORT}.

In case you want to bring down the containers the use:
```
$ docker compose -f docker/docker-compose.yml --env-file docker/.env down
```

## <img src="https://helm.sh/img/helm.svg" alt="MarineGEO circle logo" style="height: 20px; width:20px;"/> Helm

Pull the Helm Chart from the GitHub repo:

```
$ wget https://raffarus.github.io/dynamocker/helm/
```
modify the `values.yaml` file with the env variables you choose. 

Start using the UI at  http://localhost:{FE_PORT}.

## Reach the mock APIs

In order to reach the mock apis you created, you can find them at http://localhost:{BE_PORT}/dynamocker/api/serve-mock-api/<your_mock_api_url>:

```
curl http://localhost:{BE_PORT}/dynamocker/api/serve-mock-api/<your_mock_api_url>
```