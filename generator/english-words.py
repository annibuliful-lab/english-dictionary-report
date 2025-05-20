
import random
from faker import Faker
from pathlib import Path

fake = Faker()
Faker.seed(42)


def generate_realistic_word():
    choice = random.randint(1, 5)
    if choice == 1:
        return fake.word()
    elif choice == 2:
        return fake.first_name().lower()
    elif choice == 3:
        return fake.last_name().lower()
    elif choice == 4:
        return fake.color().lower()
    else:
        return fake.job().split()[0].lower()


# Output path for 500k words
file_path = Path("../data/20k.txt")


# chunk_size = 1_000_000
total_words = 20_000

with file_path.open("w") as f:
    words = [generate_realistic_word() for _ in range(total_words)]
    f.write("\n".join(words) + "\n")

file_path
