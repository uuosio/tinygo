#include <eosio/eosio.hpp>
#include <eosio/crypto.hpp>
#include <eosio/binary_extension.hpp>

using namespace std;
struct test_struct {
    string a;
    eosio::binary_extension<eosio::checksum256> b;
};

extern "C" void apply(uint64_t receiver, uint64_t first_receiver, uint64_t action) {
    if (action == "testext"_n.value) {
        auto args = eosio::unpack_action_data<test_struct>();
        eosio::check(args.b.has_value(), "no args");
        auto& value = args.b.value();
        eosio::print("value: ", value, "\n");
    } else if (action == "testext2"_n.value) {
        auto args = eosio::unpack_action_data<test_struct>();
        eosio::check(!args.b.has_value(), "no args");
    }
}
