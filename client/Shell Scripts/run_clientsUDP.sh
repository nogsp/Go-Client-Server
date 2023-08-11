# Número de clientes a serem executados
num_clients=80

# Loop para executar os clientes
for ((i=1; i<=$num_clients; i++))
do
    ./clientUDP &
done

# Aguarda todos os clientes terminarem
wait
