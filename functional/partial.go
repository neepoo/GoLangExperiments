package functional

//From a high-level point of view, our application should support dogs of multiple breeds. The breeds
//should be easily extensible. We also want to record the gender of the dog and give the dog a name.
//In our example, imagine that you’d want to spawn many dogs, so there would be a lot of repetition of
//types and genders. We’ll leverage partial application to prevent the repetitiveness of those function
//calls and improve the code readability.

const (
	Bulldog Breed = iota
	Havanese
	Cavalier
	Poodle
)
const (
	Male Gender = iota
	Female
)

var (
	maleHavaneseSpawner = DogSpawner(Havanese, Male)
	femalePoodleSpawner = DogSpawner(Poodle, Female)
)

type Dog struct {
	Name   Name
	Breed  Breed
	Gender Gender
}

type (
	Name          string
	Breed         int
	Gender        int
	NameToDogFunc func(Name) Dog
)

func DogSpawner(breed Breed, gender Gender) NameToDogFunc {
	return func(name Name) Dog {
		return Dog{
			Name:   name,
			Breed:  breed,
			Gender: gender,
		}
	}
}
