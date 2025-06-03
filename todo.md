[] Handle ping pong requests
[] Handle when FIN=0 so you gotta read and concat msgs till FIN=1
[] Handle errors with opcodes, and wrong headers and lengths
[] broadcaster resending msg to the original sender (redundant)
[] what if a client spams connect, cant keep all conn in memory need to discard the old one
[] also need to keep server side track of states, since any new joining player have to have initial values of everyone already present
[] cleanup code