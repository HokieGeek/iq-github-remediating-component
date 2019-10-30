#! /bin/bash

repo=$1
mr=$2

repo_id=$(curl -q --header "PRIVATE-TOKEN: $(cat gltoken)" \
    https://gitlab.com/api/v4/users/HokieGeek/projects | jq '.[] | select(.name=="'$repo'") | .id')

curl --header "PRIVATE-TOKEN: $(cat gltoken)" \
     https://gitlab.com/api/v4/projects/$repo_id/merge_requests/$mr/discussions \
     | jq  -c '.[] |[.id,.notes[].id]' \
     | sed 's/[]["]//g' \
     | while read l; do
        oldIFS=$IFS
        IFS=',' read -ra ids <<< "$l"
        IFS=$oldIFS
        did=${ids[0]}
        nid=${ids[1]}

        curl -L -X DELETE --header "PRIVATE-TOKEN: $(cat gltoken)" \
            https://gitlab.com/api/v4//projects/${repo_id}/merge_requests/${mr}/discussions/${did}/notes/${nid}
 done
