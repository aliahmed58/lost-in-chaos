[] Handle ping pong requests
[] Handle when FIN=0 so you gotta read and concat msgs till FIN=1
[] Handle errors with opcodes, and wrong headers and lengths
[] broadcaster resending msg to the original sender (redundant)