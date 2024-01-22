## Apis

- start working
    POST /jobs/{worktype}
        - monID
        - details depending on job type(bait, ingredients,...)

        return jobID
    DELETE /jobs/{worktype}/{jobID}
        -
    GET /jobs
    GET /jobs/{worktype}/{jobID}




## Db Migrations
Create new migrations
migrate create -dir db/migrations -ext sql -seq -digits 4 {name}




## Idea

- Collect/Hatch Monster
- monsters execute idle tasks
    - monsters have talents
        - Strength
        - Constitution
        - Intelligence
        - Wisdom

    - monsters have 4 types
     - fire
     - water
     - wind
     - earth

## Raw resources
- Woodcutting (Strength)
- Mining(Constitution)
- Farming(Inteligence)
- Fishing(Wisdom)

## PreProcessing
- Smelting(Strength)
- Carpeting
- FoodProcessing



## Finished goods
- Cooking(Wisdom, Intelligence)
- Smithing (Constitution)
- Alchemy (Health Potions, ...) (Wisdom, Intelligence)

- Hatching(make new monsters out of eggs) (depends on type) (need breeders depending on type)




- Fighting(special)

## Experience
- Collect exp by doing a job(monster get stronger)


## Traveling
- different places, which allow collecting of different stuff
- flying mounts help shorten time(endurance)
- non global market stuff can still be bought but with a delivery fee

## Future
- battles for cities
- market place
- chat based on matrix




## Next steps
- multiple jobs
- complex jobs with item consumption and cancellation
- integration testing
- master data in json files
- elements to monsters and talents
