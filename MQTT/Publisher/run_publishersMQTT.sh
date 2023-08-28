# Número de clientes a serem executados
num_clients=1

# Loop para executar os clientes
for ((i=1; i<=$num_clients; i++))
do
    ./publisher &
done

# Aguarda todos os clientes terminarem
wait
