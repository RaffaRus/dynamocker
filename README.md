## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

    helm repo add dynamocker-charts-repo https://raffarus.github.io/helm-charts

To install the dynamocker chart:

    helm install dynamocker dynamocker-charts-repo/dynamocker

To uninstall the chart:

    helm delete dynamocker
