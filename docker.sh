    #!/bin/bash
    docker build --tag forum .
    docker run --publish 6969:6969 --detach --name forumx forum